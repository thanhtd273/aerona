package com.thanhtd.aerona.booking.service;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.dto.ContactInfo;
import com.thanhtd.aerona.booking.model.Contact;

import java.util.List;

public interface ContactService {
    List<Contact> getAllContacts();

    Contact findByContactId(String contactId) throws LogicException;

    ContactInfo getContactInfo(String contactId) throws LogicException;

    ContactInfo getContactInfo(Contact contact) throws LogicException;

    Contact create(String bookingId, ContactInfo contactInfo) throws LogicException;

    Contact findByBookingId(String bookingId) throws LogicException;

    ContactInfo getBookingContact(String bookingId) throws LogicException;
}
