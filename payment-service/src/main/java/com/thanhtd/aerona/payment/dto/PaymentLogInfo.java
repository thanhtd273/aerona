package com.thanhtd.aerona.payment.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PaymentLogInfo {
    private String paymentId;

    private String status;

    private Date eventTime;

    private String description;
}
