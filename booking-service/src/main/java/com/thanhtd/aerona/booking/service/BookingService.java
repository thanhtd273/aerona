package com.thanhtd.aerona.booking.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.dto.BookingInfo;
import com.thanhtd.aerona.booking.dto.ContactInfo;
import com.thanhtd.aerona.booking.model.Booking;

import java.util.List;

public interface BookingService {

    List<Booking> getAllBookings();

    Booking findByBookingId(String bookingId) throws LogicException;

    /**
     *
     * @param bookingInfo: required field {flightId, numOfPassengers}
     * @return Booking
     * @throws LogicException
     */
    Booking initialize(BookingInfo bookingInfo) throws LogicException;

    Booking book(BookingInfo bookingInfo) throws LogicException, JsonProcessingException;

    BookingInfo getBookingDetail(String bookingId) throws LogicException;

    BookingInfo getBookingInfo(Booking booking) throws LogicException;

    ContactInfo getBookingContact(String bookingId) throws LogicException;

    Booking updateBookingStatus(String bookingId, String status) throws LogicException;

}
