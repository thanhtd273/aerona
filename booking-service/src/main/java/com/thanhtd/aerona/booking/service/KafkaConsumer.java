package com.thanhtd.aerona.booking.service;

public interface KafkaConsumer {

    void handleUpdateStatus(String bookingPayload);
}
