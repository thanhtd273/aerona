package com.thanhtd.aerona.booking.constants;

public final class TopicConstant {

    public static final String BOOKING_STATUS_UPDATED = "booking.updated-status";

    public static final String PAYMENT_COMPLETED = "payment.completed";

    private TopicConstant() {
        throw new IllegalStateException("Utility class");
    }
}
