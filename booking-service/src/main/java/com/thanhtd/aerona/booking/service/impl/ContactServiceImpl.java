package com.thanhtd.aerona.booking.service.impl;

import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.AppUtils;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.dao.ContactDao;
import com.thanhtd.aerona.booking.dto.ContactInfo;
import com.thanhtd.aerona.booking.model.Contact;
import com.thanhtd.aerona.booking.service.ContactService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.util.Date;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ContactServiceImpl implements ContactService {

    private final ContactDao contactDao;

    @Override
    public List<Contact> getAllContacts() {
        return contactDao.findAll();
    }

    @Override
    public Contact findByContactId(String contactId) throws LogicException {
        if (ObjectUtils.isEmpty(contactId)) {
            throw new LogicException(ErrorCode.ID_NULL);
        }
        return contactDao.findByContactId(contactId);
    }

    @Override
    public ContactInfo getContactInfo(String contactId) throws LogicException {
        Contact contact = findByContactId(contactId);
        return new ContactInfo(contact.getFirstName(), contact.getLastName(), contact.getPhone(), contact.getEmail());
    }

    @Override
    public ContactInfo getContactInfo(Contact contact) throws LogicException {
        if (ObjectUtils.isEmpty(contact)) {
            throw new LogicException(ErrorCode.DATA_NULL, "Not found contact");
        }
        return new ContactInfo(contact.getFirstName(), contact.getLastName(), contact.getPhone(), contact.getEmail());
    }

    @Override
    public Contact create(String bookingId, ContactInfo contactInfo) throws LogicException {
        if (ObjectUtils.isEmpty(contactInfo) || ObjectUtils.isEmpty(bookingId)) {
            throw new LogicException(ErrorCode.DATA_NULL);
        }
        if (!contactInfo.areAllFieldsFilled()) {
            throw new LogicException(ErrorCode.BLANK_FIELD);
        }
        Contact contact = new Contact();
        contact.setContactId(UUID.randomUUID().toString());
        contact.setBookingId(bookingId);
        contact.setFirstName(contactInfo.getFirstName());
        contact.setLastName(contactInfo.getLastName());

        if (!AppUtils.validateEmail(contactInfo.getEmail()))
            throw new LogicException(ErrorCode.INVALID_DATA_FORMAT, "Invalid email");
        contact.setEmail(contactInfo.getEmail());

        if (!AppUtils.validatePhoneNumber(contactInfo.getPhone()))
            throw new LogicException(ErrorCode.INVALID_DATA_FORMAT, "Phone number must start with 0 and have 10 digits");
        contact.setPhone(contactInfo.getPhone());

        contact.setCreatedAt(new Date(System.currentTimeMillis()));
        contact.setStatus(DataStatus.ACTIVE);
        return contactDao.save(contact);
    }

    @Override
    public Contact findByBookingId(String bookingId) throws LogicException {
        if (ObjectUtils.isEmpty(bookingId)) {
            throw new LogicException(ErrorCode.ID_NULL);
        }
        return contactDao.findByBookingId(bookingId);
    }

    @Override
    public ContactInfo getBookingContact(String bookingId) throws LogicException {
        Contact contact = findByBookingId(bookingId);
        return getContactInfo(contact);
    }
}
