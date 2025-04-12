package com.thanhtd.aerona.booking.service.impl;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.constants.TopicConstant;
import com.thanhtd.aerona.booking.dto.BookingInfo;
import com.thanhtd.aerona.booking.constants.GroupConstant;
import com.thanhtd.aerona.booking.service.BookingService;
import com.thanhtd.aerona.booking.service.KafkaConsumer;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@AllArgsConstructor
@Slf4j
public class KafkaConsumerImpl implements KafkaConsumer {


    private final BookingService bookingService;

    private final ObjectMapper objectMapper;

    @Override
    @KafkaListener(topics = TopicConstant.BOOKING_STATUS_UPDATED, groupId = GroupConstant.FLIGHT_GROUP)
    public void handleUpdateStatus(String bookingPayload) {
        BookingInfo booking = new BookingInfo();
        try {
            booking = objectMapper.readValue(bookingPayload, BookingInfo.class);
            bookingService.updateBookingStatus(booking.getBookingId(), booking.getStatus());
            log.debug("Update status of booking with booking_id={} to {}", booking.getBookingId(), booking.getStatus());
        } catch (JsonProcessingException e) {
            log.error("Failed to convert JSON to Java object");
        } catch(LogicException e) {
            log.error("Failed to update status of booking {} to {}, error: {}", booking.getBookingId(), booking.getStatus(), e.getMessage());
        }
    }


}
