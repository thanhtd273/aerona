package com.thanhtd.aerona.payment.constant;

public final class TopicConstant {

    public static final String PAYMENT_CREATED = "payment.created";

    public static final String PAYMENT_COMPLETED = "payment.completed";

    public static final String PAYMENT_FAILED = "payment.failed";

    private TopicConstant() {
        throw new IllegalStateException("Utility class");
    }
}
