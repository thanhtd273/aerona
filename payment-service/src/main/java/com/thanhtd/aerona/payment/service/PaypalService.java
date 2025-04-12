package com.thanhtd.aerona.payment.service;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.payment.dto.PaymentInfo;
import com.thanhtd.aerona.payment.dto.PaymentSource;
import com.thanhtd.aerona.payment.dto.paypal.OrderResponse;

public interface PaypalService {

    String getAccessToken() throws LogicException;

    OrderResponse createOrder(PaymentInfo paymentInfo) throws LogicException;

    OrderResponse confirmPaymentSource(String orderId, PaymentSource paymentSource) throws LogicException;

    OrderResponse captureOrder(String orderId) throws LogicException;
}
