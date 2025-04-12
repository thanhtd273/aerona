package com.thanhtd.aerona.booking.service.impl;

import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.dao.PassengerDao;
import com.thanhtd.aerona.booking.dto.PassengerInfo;
import com.thanhtd.aerona.booking.model.Passenger;
import com.thanhtd.aerona.booking.service.PassengerService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.util.Date;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class PassengerServiceImpl implements PassengerService {

    private final PassengerDao passengerDao;

    @Override
    public Passenger findByPassengerId(String passengerId) throws LogicException {
        if (ObjectUtils.isEmpty(passengerId)) {
            throw new LogicException(ErrorCode.ID_NULL);
        }

        return passengerDao.findByPassengerId(passengerId);
    }

    @Override
    public List<PassengerInfo> getInfoByBookingId(String bookingId) throws LogicException {
        List<Passenger> passengers = findByBookingId(bookingId);
        return passengers.stream().map(p -> new PassengerInfo(p.getFirstName(), p.getLastName(),
                p.getDayOfBirth(), p.getNationality(), p.getPassportNumber())).toList();
    }

    @Override
    public List<Passenger> findByBookingId(String bookingId) throws LogicException {
        if (ObjectUtils.isEmpty(bookingId)) {
            throw new LogicException(ErrorCode.ID_NULL);
        }
        return passengerDao.findByBookingId(bookingId);
    }

    @Override
    public PassengerInfo getPassengerInfo(String passengerId) throws LogicException {
        Passenger passenger = findByPassengerId(passengerId);
        return new PassengerInfo(passenger.getFirstName(), passenger.getLastName(),
                passenger.getDayOfBirth(), passenger.getNationality(), passenger.getPassportNumber());
    }

    @Override
    public Passenger create(String bookingId, PassengerInfo passengerInfo) throws LogicException {
        if (ObjectUtils.isEmpty(bookingId) || ObjectUtils.isEmpty(passengerInfo)) {
            throw new LogicException(ErrorCode.DATA_NULL);
        }
        if (!passengerInfo.areFieldsFilled()) {
            throw new LogicException(ErrorCode.BLANK_FIELD);
        }
        Passenger passenger = new Passenger();
        passenger.setPassengerId(UUID.randomUUID().toString());
        passenger.setBookingId(bookingId);
        passenger.setFirstName(passengerInfo.getFirstName());
        passenger.setLastName(passengerInfo.getLastName());
        passenger.setNationality(passengerInfo.getNationality());
        passenger.setDayOfBirth(passengerInfo.getDayOfBirth());
        passenger.setPassportNumber(passengerInfo.getPassportNumber());
        passenger.setCreatedAt(new Date(System.currentTimeMillis()));
        passenger.setStatus(DataStatus.ACTIVE);
        return passengerDao.save(passenger);
    }

    @Override
    public void createMany(String bookingId, List<PassengerInfo> passengerInfos) throws LogicException {
        for (PassengerInfo passengerInfo : passengerInfos) {
            create(bookingId, passengerInfo);
        }
    }
}
