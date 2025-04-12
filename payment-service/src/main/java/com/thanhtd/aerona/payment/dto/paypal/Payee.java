package com.thanhtd.aerona.payment.dto.paypal;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Setter
@Getter
public class Payee {

    @JsonProperty("email_address")
    private String emailAddress;

    @JsonProperty("merchant_id")
    private String merchantId;
}
