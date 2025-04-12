package com.thanhtd.aerona.payment.constant;

public final class BookingStatus {
    
    public static final String INITIAL = "INITIAL";

    public  static final String BOOK = "BOOK";

    public static final String PENDING = "PENDING";

    public static final String CANCELED = "CANCELED";

    public static final String COMPLETED = "COMPLETED";

    public static final String PAYMENT_FAILED = "PAYMENT_FAILED";

    public static final String TICKET_ISSUED = "TICKET_ISSUED";

    public static final String TICKET_FAILED = "TICKET_FAILED";

    public static final String NOTIFIED = "NOTIFIED";

    public static final String NOTIFICATION_FAILED = "NOTIFICATION_FAILED";


    private BookingStatus() {
        throw new IllegalStateException("Utility class");
    }
}
