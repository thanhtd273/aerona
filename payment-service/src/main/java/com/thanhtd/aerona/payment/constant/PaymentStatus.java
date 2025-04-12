package com.thanhtd.aerona.payment.constant;

public final class PaymentStatus {

    public static final String DONE = "DONE";

    public static final String FAILED = "FAILED";

    private PaymentStatus() {
        throw new IllegalStateException("Utility class");
    }
}
