package com.thanhtd.aerona.booking.dao;

import com.thanhtd.aerona.booking.model.Contact;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ContactDao extends JpaRepository<Contact, String> {

    Contact findByContactId(String contactId);

    Contact findByBookingId(String bookingId);
}
