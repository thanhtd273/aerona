package com.thanhtd.aerona.payment.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

import java.util.Date;

@Entity
@Table(name = "payment")
@Data
public class Payment {
    @Id
    @Column(name = "payment_id")
    private String paymentId;

    @Column(name = "booking_id")
    private String bookingId;

    @Column(name = "pnr")
    private String pnr;

    @Column(name = "payment_method")
    private String paymentMethod;

    @Column(name = "payment_platform")
    private String paymentPlatform;

    @Column(name = "external_payment_id")
    private String externalPaymentId;

    @Column(name = "amount")
    private Float amount;

    @Column(name = "currency")
    private String currency;

    @Column(name = "status")
    private String status;

    @Column(name = "created_at")
    private Date createdAt;

    @Column(name = "paid_at")
    private Date paidAt;

    @Column(name = "failed_at")
    private Date failedAt;

    @Column(name = "refunded_at")
    private Date refundedAt;

    @Column(name = "failure_reason")
    private String failureReason;
}
