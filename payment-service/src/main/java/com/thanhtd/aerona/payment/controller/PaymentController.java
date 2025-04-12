package com.thanhtd.aerona.payment.controller;

import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import com.thanhtd.aerona.payment.dto.PaymentInfo;
import com.thanhtd.aerona.payment.dto.PaymentSource;
import com.thanhtd.aerona.payment.model.Payment;
import com.thanhtd.aerona.payment.service.PaymentService;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping(value = "/api/v1/payments")
@RequiredArgsConstructor
@Slf4j
public class PaymentController {

    private final PaymentService paymentService;

    @GetMapping()
    public APIResponse getAllPayments(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Payment> payments = paymentService.getAllPayments();
            log.debug("Call GET /api/v1/payments success, count={}", payments.size());
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, payments);
        } catch (Exception e) {
            log.debug("Call GET /api/v1/payments failed. error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping()
    public APIResponse createPayment(HttpServletResponse response, @RequestBody PaymentInfo paymentInfo) {
        long start = System.currentTimeMillis();
        try {
            Payment payment = paymentService.createPaymentThroughPaypal(paymentInfo);
            log.debug("Call POST /api/v1/payments success, paymentId={}", payment.getPaymentId());
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, payment);
        } catch (Exception e) {
            log.debug("Call POST /api/v1/payments failed. error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/{paymentId}/confirm")
    public APIResponse confirmPayment(@PathVariable("paymentId") String paymentId, @RequestBody PaymentSource paymentSource,
                                      HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            Payment payment = paymentService.confirmPayment(paymentId, paymentSource);
            log.debug("Call POST /api/v1/payments/{}/confirm success, status: {}", payment.getPaymentId(), payment.getStatus());
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, payment);
        } catch (Exception e) {
            log.error("Call POST /api/v1/payments/{}/confirm failed, error: {}", paymentId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/{paymentId}/capture")
    public APIResponse confirmPayment(@PathVariable("paymentId") String paymentId, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            Payment payment = paymentService.capturePayment(paymentId);
            log.debug("Call POST /api/v1/payments/{}/capture success, status: {}", payment.getPaymentId(), payment.getStatus());
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, payment);
        } catch (Exception e) {
            log.error("Call POST /api/v1/payments/{}/capture failed, error: {}", paymentId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }
}
