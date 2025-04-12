package com.thanhtd.aerona.payment.dao;

import com.thanhtd.aerona.payment.model.Refund;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface RefundDao extends JpaRepository<Refund, String> {
}
