package com.thanhtd.aerona.booking.model;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;

@Data
@Entity
@Table(name = "passenger")
public class Passenger {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "passenger_id")
    private String passengerId;

    @Column(name = "booking_id")
    private String bookingId;

    @Column(name = "first_name")
    private String firstName;

    @Column(name = "last_name")
    private String lastName;

    @Column(name = "day_of_birth")
    private Date dayOfBirth;

    @Column(name = "nationality")
    private String nationality;

    @Column(name = "passport_number")
    private String passportNumber;

    @Column(name = "created_at")
    private Date createdAt;

    @Column(name = "updated_at")
    private Date updatedAt;

    @Column(name = "status")
    private String status;
}
