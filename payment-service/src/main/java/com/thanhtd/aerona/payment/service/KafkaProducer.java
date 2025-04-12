package com.thanhtd.aerona.payment.service;

import com.thanhtd.aerona.payment.dto.BookingInfo;

public interface KafkaProducer {

    void sendPaymentCreated(BookingInfo paymentInfo);

    void sendPaymentCompleted(BookingInfo paymentInfo);

    void sendPaymentFailed(BookingInfo paymentInfo);
}
