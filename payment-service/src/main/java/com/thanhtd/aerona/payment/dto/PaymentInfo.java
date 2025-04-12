package com.thanhtd.aerona.payment.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PaymentInfo {

    @NotBlank(message = "bookingId is not empty")
    @JsonProperty("booking_id")
    private String bookingId;

    @NotNull(message = "PRN code is not empty")
    private String pnr;

    @NotNull(message = "amount is not empty")
    @Min(value = 1, message = "amount must be greater than 0")
    private Float amount;

    @NotBlank(message = "currency is not empty")
    private String currency;

    private String status;
}
