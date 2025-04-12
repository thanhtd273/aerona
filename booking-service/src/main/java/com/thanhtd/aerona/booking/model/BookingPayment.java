package com.thanhtd.aerona.booking.model;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;

@Data
@Entity
@Table(name = "booking_payment")
public class BookingPayment {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "booking_id")
    private String bookingId;

    @Column(name = "payment_id")
    private String paymentId;

    @Column(name = "amount")
    private Float amount;

    @Column(name = "paid_at")
    private Date paidAt;
}
