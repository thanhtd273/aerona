package com.thanhtd.aerona.payment.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

import java.util.Date;

@Entity
@Table(name = "payment_log")
@Data
public class PaymentLog {
    @Id
    @Column(name = "id")
    private String id;

    @Column(name = "payment_id")
    private String paymentId;

    @Column(name = "status")
    private String status;

    @Column(name = "event_time")
    private Date eventTime;

    @Column(name = "description")
    private String description;
}
