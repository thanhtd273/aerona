package com.thanhtd.aerona.booking.dao;

import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.booking.model.Booking;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface BookingDao extends JpaRepository<Booking, String> {
    Booking findByBookingId(String bookingId) throws LogicException;

    Integer countByPnr(String pnr);
}
