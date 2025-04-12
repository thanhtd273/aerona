package com.thanhtd.aerona.booking.service;

import com.thanhtd.aerona.booking.dto.BookingInfo;

public interface KafkaProducer {

    void sendPaymentCompleted(BookingInfo bookingInfo);

}
