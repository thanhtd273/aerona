package com.thanhtd.aerona.payment.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

import java.util.Date;

@Entity
@Table(name = "refund")
@Data
public class Refund {
    @Id
    @Column(name = "refund_id")
    private String refundId;

    @Column(name = "payment_id")
    private String paymentId;

    @Column(name = "external_refund_id")
    private String externalRefundId;

    @Column(name = "amount")
    private Float amount;

    @Column(name = "currency")
    private String currency;

    @Column(name = "requested_at")
    private Date requestedAt;

    @Column(name = "completed_at")
    private Date completedAt;

    @Column(name = "status")
    private String status;

    @Column(name = "reason")
    private String reason;
}
