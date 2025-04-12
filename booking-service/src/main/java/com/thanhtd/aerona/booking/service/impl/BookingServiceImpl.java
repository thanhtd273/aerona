package com.thanhtd.aerona.booking.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.booking.dto.BookingInfo;
import com.thanhtd.aerona.booking.dto.ContactInfo;
import com.thanhtd.aerona.booking.dto.Flight;
import com.thanhtd.aerona.booking.dto.PassengerInfo;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.booking.dto.*;
import com.thanhtd.aerona.booking.service.ContactService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import org.springframework.beans.factory.annotation.Value;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.constants.BookingStatus;
import com.thanhtd.aerona.booking.dao.BookingDao;
import com.thanhtd.aerona.booking.model.Booking;
import com.thanhtd.aerona.booking.model.Contact;
import com.thanhtd.aerona.booking.service.BookingService;
import com.thanhtd.aerona.booking.service.PassengerService;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.ObjectUtils;
import org.springframework.web.client.RestTemplate;

import java.util.Date;
import java.util.List;
import java.util.Objects;
import java.util.Random;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class BookingServiceImpl implements BookingService {

    private static final int EXPIRE_TIME = 30 * 60 * 1000;

    private static final String CHARACTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

    private static final int PNR_LENGTH = 6;

    private static final Random RANDOM = new Random();

    @Value("${flight.api.base-url}")
    private String baseUrl;

    private final BookingDao bookingDao;

    private final ContactService contactService;

    private final PassengerService passengerService;

    private final RestTemplate restTemplate;

    private final ObjectMapper objectMapper;


    @Override
    public List<Booking> getAllBookings() {
        return bookingDao.findAll();
    }

    @Override
    public Booking findByBookingId(String bookingId) throws LogicException {
        if (ObjectUtils.isEmpty(bookingId)) {
            throw new LogicException(ErrorCode.ID_NULL);
        }
        return bookingDao.findByBookingId(bookingId);
    }

    @Override
    public Booking initialize(BookingInfo bookingInfo) throws LogicException {
        if (ObjectUtils.isEmpty(bookingInfo)) {
            throw new LogicException(ErrorCode.DATA_NULL);
        }
        String flightId = bookingInfo.getFlightId();
        Integer numOfPassengers = bookingInfo.getNumOfPassengers();
        if (ObjectUtils.isEmpty(flightId) || numOfPassengers <= 0) {
            throw new LogicException(ErrorCode.ID_NULL, "Flight's ID is empty or passenger is invalid");
        }
        Flight flight = findFlightById(flightId);
        if (ObjectUtils.isEmpty(flight)) {
            throw new LogicException(ErrorCode.DATA_NULL, "Flight does not exist");
        }

        Booking booking = new Booking();
        booking.setBookingId(UUID.randomUUID().toString());
        String pnr = generateUniquePNR();
        booking.setPnr(pnr);
        booking.setFlightId(bookingInfo.getFlightId());
        if (!ObjectUtils.isEmpty(bookingInfo.getUserId()))
            booking.setUserId(bookingInfo.getUserId());

        booking.setNumOfPassengers(bookingInfo.getNumOfPassengers());
        booking.setTotalPrice(flight.getPrice() * numOfPassengers);
        booking.setCurrency(flight.getCurrency());
        booking.setBookedAt(new Date(System.currentTimeMillis()));
        booking.setStatus(BookingStatus.INITIAL);
        booking.setExpiresAt(new Date(System.currentTimeMillis() + EXPIRE_TIME));

        return bookingDao.save(booking);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public Booking book(BookingInfo bookingInfo) throws LogicException, JsonProcessingException {
        if (ObjectUtils.isEmpty(bookingInfo)) {
            throw new LogicException(ErrorCode.DATA_NULL);
        }
        if (ObjectUtils.isEmpty(bookingInfo.getContact()) || ObjectUtils.isEmpty(bookingInfo.getPassengers())) {
            throw new LogicException(ErrorCode.DATA_NULL, "Contact or passengers who belongs to booking is empty");
        }
        String bookingId = bookingInfo.getBookingId();
        Booking booking = findByBookingId(bookingId);
        if (ObjectUtils.isEmpty(booking)) {
            throw new LogicException(ErrorCode.DATA_NULL, "Not found booking with id " + bookingId);
        }
        if (!Objects.equals(bookingInfo.getPassengers().size(), booking.getNumOfPassengers())) {
            throw new LogicException(ErrorCode.DATA_NULL, "Number of passengers does not match");
        }
        Contact contact = contactService.create(bookingId, bookingInfo.getContact());
        passengerService.createMany(bookingId, bookingInfo.getPassengers());
        booking.setContactId(contact.getContactId());
        booking.setStatus(BookingStatus.BOOK);
        booking.setUpdatedAt(new Date(System.currentTimeMillis()));

        booking = bookingDao.save(booking);

        return booking;
    }

    @Override
    public BookingInfo getBookingDetail(String bookingId) throws LogicException {
        Booking booking = findByBookingId(bookingId);
        return getBookingInfo(booking);
    }


    @Override
    public BookingInfo getBookingInfo(Booking booking) throws LogicException {
        if (ObjectUtils.isEmpty(booking)) {
            throw new LogicException(ErrorCode.DATA_NULL, "Booking is null");
        }
        ContactInfo contact = contactService.getContactInfo(booking.getContactId());
        List<PassengerInfo> passengers = passengerService.getInfoByBookingId(booking.getBookingId());
        Flight flight = findFlightById(booking.getFlightId());
        return new BookingInfo(booking.getBookingId(), booking.getPnr(), booking.getFlightId(), flight, booking.getNumOfPassengers(),
                booking.getUserId(), contact, passengers, booking.getTotalPrice(), booking.getCurrency(), booking.getStatus());
    }

    @Override
    public ContactInfo getBookingContact(String bookingId) throws LogicException {
        return contactService.getBookingContact(bookingId);
    }


    @Override
    public Booking updateBookingStatus(String bookingId, String status) throws LogicException {
        Booking booking = findByBookingId(bookingId);
        if (ObjectUtils.isEmpty(booking)) {
            throw new LogicException(ErrorCode.DATA_NULL, "Booking does not exist");
        }
        booking.setStatus(status);
        return bookingDao.save(booking);
    }

    private Flight findFlightById(String flightId) {
        String url = baseUrl + "/flights/" + flightId;
        APIResponse response = restTemplate.getForObject(url, APIResponse.class);
        if (ObjectUtils.isEmpty(response)) {
            return null;
        }
        return objectMapper.convertValue(response.getData(), Flight.class);
    }

    private String generateUniquePNR() throws LogicException {
        String pnr;
        int maxAttempt = 10;
        for (int attempt = 0; attempt < maxAttempt; attempt++) {
            pnr = generateRandomPNR();
            if (isUniquePNR(pnr))
                return pnr;
        }
        throw new LogicException(ErrorCode.DATA_NULL, String.format("Unable to generate unique PNR after %d attempts", maxAttempt));
    }

    private String generateRandomPNR() {
        StringBuilder pnr = new StringBuilder();
        for (int i = 0; i < PNR_LENGTH; i++) {
            int index = RANDOM.nextInt(CHARACTERS.length());
            pnr.append(CHARACTERS.charAt(index));
        }
        return pnr.toString();
    }

    private boolean isUniquePNR(String pnr) {
        Integer count = bookingDao.countByPnr(pnr);
        return count != null && count == 0;
    }


}
