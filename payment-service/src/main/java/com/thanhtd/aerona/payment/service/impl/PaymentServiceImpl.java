package com.thanhtd.aerona.payment.service.impl;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.core.AppUtils;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.payment.dao.PaymentDao;
import com.thanhtd.aerona.payment.dto.BookingInfo;
import com.thanhtd.aerona.payment.dto.PaymentInfo;
import com.thanhtd.aerona.payment.dto.PaymentSource;
import com.thanhtd.aerona.payment.dto.paypal.OrderResponse;
import com.thanhtd.aerona.payment.model.Payment;
import com.thanhtd.aerona.payment.service.KafkaProducer;
import com.thanhtd.aerona.payment.service.PaymentService;

import com.thanhtd.aerona.payment.service.PaypalService;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.Validator;
import jakarta.validation.ValidatorFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;
import org.springframework.web.client.RestTemplate;

import java.util.Date;
import java.util.List;
import java.util.Set;

@Service
@RequiredArgsConstructor
@Slf4j
public class PaymentServiceImpl implements PaymentService {

    @Value("${booking.service.base.url}")
    private String bookingServiceUrl;

    private final PaymentDao paymentDao;

    private final PaypalService paypalService;

    private final KafkaProducer kafkaProducer;

    private final RestTemplate restTemplate;

    private final ObjectMapper objectMapper;

    @Override
    public List<Payment> getAllPayments() {
        return paymentDao.findAll();
    }

    @Override
    public Payment findByPaymentId(String paymentId) {
        return paymentDao.findByPaymentId(paymentId);
    }

    @Override
    public Payment findByBookingId(String bookingId) {
        return paymentDao.findByBookingId(bookingId);
    }

    @Override
    public Payment createPaymentThroughPaypal(PaymentInfo paymentInfo) throws LogicException {
        if (ObjectUtils.isEmpty(paymentInfo)) {
            throw new LogicException(ErrorCode.DATA_NULL);
        }
        Validator validator;
        try (ValidatorFactory factory = Validation.buildDefaultValidatorFactory()) {
            validator = factory.getValidator();
        }
        Set<ConstraintViolation<PaymentInfo>> violations = validator.validate(paymentInfo);
        if (!violations.isEmpty()) {
            violations.forEach(v -> log.error("Invalid field: {}", v.getMessage()));
            throw new LogicException(ErrorCode.BLANK_FIELD);
        }
        OrderResponse response = paypalService.createOrder(paymentInfo);

        Payment payment = new Payment();
        payment.setExternalPaymentId(response.getId());
        payment.setPaymentId(AppUtils.generateUniqueId());
        payment.setBookingId(paymentInfo.getBookingId());
        payment.setPnr(paymentInfo.getPnr());
        payment.setAmount(paymentInfo.getAmount());
        payment.setCurrency(paymentInfo.getCurrency());
        payment.setPaymentMethod("CREDIT_CARD");
        payment.setPaymentPlatform("PAYPAL");
        payment.setCreatedAt(new Date(System.currentTimeMillis()));
        payment.setStatus(response.getStatus());

        BookingInfo bookingInfo = getBookingDetail(paymentInfo.getBookingId());

       kafkaProducer.sendPaymentCreated(bookingInfo);
        return paymentDao.save(payment);
    }

    @Override
    public Payment confirmPayment(String paymentId, PaymentSource paymentSource) throws LogicException {
        Payment payment = findByPaymentId(paymentId);
        if (ObjectUtils.isEmpty(payment)) {
            throw new LogicException(ErrorCode.NOT_FOUND_PAYMENT);
        }
        OrderResponse response = paypalService.confirmPaymentSource(payment.getExternalPaymentId(), paymentSource);
        payment.setStatus(response.getStatus());
        paymentDao.save(payment);
        return payment;
    }

    @Override
    public Payment capturePayment(String paymentId) throws LogicException {
        Payment payment = findByPaymentId(paymentId);
        if (ObjectUtils.isEmpty(payment)) {
            throw new LogicException(ErrorCode.NOT_FOUND_PAYMENT);
        }
        OrderResponse response = paypalService.captureOrder(payment.getExternalPaymentId());
        payment.setPaidAt(new Date(System.currentTimeMillis()));
        payment.setStatus(response.getStatus());
        paymentDao.save(payment);
        BookingInfo bookingInfo = getBookingDetail(payment.getBookingId());
        kafkaProducer.sendPaymentCompleted(bookingInfo);
        return payment;
    }

    private BookingInfo getBookingDetail(String bookingId) {
        String url = bookingServiceUrl + "/api/v1/bookings/" + bookingId + "/detail";
        APIResponse response = restTemplate.getForObject(url, APIResponse.class);
        if (ObjectUtils.isEmpty(response)) {
            return null;
        }
        return objectMapper.convertValue(response.getData(), BookingInfo.class);
    }
}
