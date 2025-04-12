package com.thanhtd.aerona.payment.service;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.payment.dto.PaymentLogInfo;
import com.thanhtd.aerona.payment.model.PaymentLog;

import java.util.List;

public interface PaymentLogService {
    List<PaymentLog> getAllPaymentLogs();

    PaymentLog findByPaymentLogId(String logId);

    List<PaymentLog> findByPaymentId(String paymentId);

    PaymentLog createPaymentLog(PaymentLogInfo paymentLogInfo) throws LogicException;
}
