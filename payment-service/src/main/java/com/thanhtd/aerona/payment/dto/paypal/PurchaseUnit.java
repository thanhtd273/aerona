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
public class PurchaseUnit {
    @JsonProperty("reference_id")
    private String referenceId;

    private PaymentAmount amount;

    private Payee payee;

    private String description;

    @JsonProperty("custom_id")
    private String customId;

    @JsonProperty("invoice_id")
    private String invoiceId;

    @JsonProperty("soft_descriptor")
    private String softDescriptor;

}
