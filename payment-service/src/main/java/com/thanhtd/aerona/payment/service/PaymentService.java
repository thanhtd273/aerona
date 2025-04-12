package com.thanhtd.aerona.payment.service;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.payment.dto.PaymentInfo;
import com.thanhtd.aerona.payment.dto.PaymentSource;
import com.thanhtd.aerona.payment.model.Payment;

import java.util.List;

public interface PaymentService {

    List<Payment> getAllPayments();

    Payment findByPaymentId(String paymentId);

    Payment findByBookingId(String bookingId);

    Payment createPaymentThroughPaypal(PaymentInfo paymentInfo) throws LogicException;

    Payment confirmPayment(String paymentId, PaymentSource paymentSource) throws LogicException;

    Payment capturePayment(String paymentId) throws LogicException;
}
