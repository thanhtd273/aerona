package com.thanhtd.aerona.payment.dao;

import com.thanhtd.aerona.payment.model.Payment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface PaymentDao extends JpaRepository<Payment, String> {

    Payment findByPaymentId(String paymentId);

    Payment findByBookingId(String bookingId);
    
}
