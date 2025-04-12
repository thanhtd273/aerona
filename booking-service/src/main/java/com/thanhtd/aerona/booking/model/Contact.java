package com.thanhtd.aerona.booking.model;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;

@Data
@Entity
@Table(name= "contact")
public class Contact {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "contact_id")
    private String contactId;

    @Column(name = "booking_id")
    private String bookingId;

    @Column(name = "first_name")
    private String firstName;

    @Column(name = "last_name")
    private String lastName;

    @Column(name = "email")
    private String email;

    @Column(name = "phone")
    private String phone;

    @Column(name = "created_at")
    private Date createdAt;

    @Column(name = "updated_at")
    private Date updatedAt;

    @Column(name = "status")
    private String status;
}
