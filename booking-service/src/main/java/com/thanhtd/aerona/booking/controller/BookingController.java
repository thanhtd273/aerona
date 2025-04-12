package com.thanhtd.aerona.booking.controller;

import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import com.thanhtd.aerona.booking.dto.BookingInfo;
import com.thanhtd.aerona.booking.dto.ContactInfo;
import com.thanhtd.aerona.booking.model.Booking;
import com.thanhtd.aerona.booking.service.BookingService;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequiredArgsConstructor
public class BookingController {
    private static final Logger logger = LoggerFactory.getLogger(BookingController.class);

    private static final String SUCCESS_DESCRIPTION = "success";

    private final BookingService bookingService;

    @GetMapping("/internal/bookings")
    public APIResponse getAllBookings(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Booking> bookings = bookingService.getAllBookings();
            logger.debug("Bookings found: {}", bookings);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, bookings);
        } catch (Exception e) {
            logger.error("Failed to call GET /internal/bookings, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping("/internal/bookings/{bookingId}")
    public APIResponse getBooking(@PathVariable String bookingId, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            Booking booking = bookingService.findByBookingId(bookingId);
            logger.debug("Booking found: {}", booking);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, booking);
        } catch(Exception e) {
            logger.error("Failed to call GET /internal/bookings/{}, error: {}", bookingId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "/api/v1/bookings/initialize")
    public APIResponse initialize(@RequestBody BookingInfo bookingInfo, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            Booking booking = bookingService.initialize(bookingInfo);
            logger.debug("Booking initialized: {}", booking);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, booking);
        } catch (Exception e) {
            logger.error("Failed to call POST /api/v1/bookings/initialize, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "/api/v1/bookings/book")
    public APIResponse book(@RequestBody BookingInfo bookingInfo, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            Booking booking = bookingService.book(bookingInfo);
            logger.debug("Booking booked: {}", booking);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, booking);
        } catch (Exception e) {
            logger.error("Failed to call POST /api/v1/bookings/book, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/api/v1/bookings/{bookingId}/detail")
    public APIResponse getBookingDetail(@PathVariable String bookingId, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            BookingInfo bookingInfo = bookingService.getBookingDetail(bookingId);
            logger.debug("Booking detail found: {}", bookingInfo);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, bookingInfo);
        } catch (Exception e) {
            logger.error("Failed to call GET /api/v1/bookings/{}/detail, error: {}", bookingId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/internal/bookings/{bookingId}/contact")
    public APIResponse getBookingContact(@PathVariable String bookingId, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            ContactInfo contactInfo = bookingService.getBookingContact(bookingId);
            logger.debug("Booking contact found: {}", contactInfo);
            return new APIResponse(ErrorCode.SUCCESS, SUCCESS_DESCRIPTION, System.currentTimeMillis() - start, contactInfo);
        } catch (Exception e) {
            logger.error("Failed to call GET /api/v1/bookings/{}/contact, error: {}", bookingId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }
}
