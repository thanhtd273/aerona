package com.thanhtd.aerona.payment.dto.paypal;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PaymentAmount {
    @JsonProperty("currency_code")
    private String currencyCode;

    private String value;
}
