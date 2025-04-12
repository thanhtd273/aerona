package com.thanhtd.aerona.payment.dao;

import com.thanhtd.aerona.payment.model.PaymentLog;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface PaymentLogDao extends JpaRepository<PaymentLog, String> {
}
