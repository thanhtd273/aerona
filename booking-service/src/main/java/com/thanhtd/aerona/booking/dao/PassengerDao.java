package com.thanhtd.aerona.booking.dao;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.model.Passenger;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PassengerDao extends JpaRepository<Passenger, String> {

    Passenger findByPassengerId(String passengerId) throws LogicException;

    List<Passenger> findByBookingId(String bookingId) throws LogicException;
}
