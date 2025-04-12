package com.thanhtd.aerona.booking.service;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.dto.PassengerInfo;
import com.thanhtd.aerona.booking.model.Passenger;

import java.util.List;

public interface PassengerService {

    List<Passenger> findByBookingId(String bookingId) throws LogicException;

    List<PassengerInfo> getInfoByBookingId(String bookingId) throws LogicException;

    PassengerInfo getPassengerInfo(String passengerId) throws LogicException;

    Passenger findByPassengerId(String passengerId) throws LogicException;

    Passenger create(String bookingId, PassengerInfo passengerInfo) throws LogicException;

    void createMany(String bookingId, List<PassengerInfo> passengerInfos) throws LogicException;
}
