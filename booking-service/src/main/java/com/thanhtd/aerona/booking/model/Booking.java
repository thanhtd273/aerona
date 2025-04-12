package com.thanhtd.aerona.booking.model;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;

@Data
@Entity
@Table(name = "booking")
public class Booking {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "booking_id")
    private String bookingId;

    @Column(name = "pnr")
    private String pnr;

    @Column(name = "flight_id")
    private String flightId;

    @Column(name = "user_id")
    private String userId;

    @Column(name = "contact_id")
    private String contactId;

    @Column(name = "num_of_passengers")
    private Integer numOfPassengers;

    @Column(name = "status")
    private String status;

    @Column(name = "booked_at")
    private Date bookedAt;

    @Column(name = "expires_at")
    private Date expiresAt;

    @Column(name = "confirmed_at")
    private Date confirmedAt;

    @Column(name = "cancelled_at")
    private Date cancelledAt;

    @Column(name = "updated_at")
    private Date updatedAt;

    @Column(name = "total_price")
    private Float totalPrice;

    @Column(name = "currency")
    private String currency;
}
