package com.thanhtd.aerona.payment.dto.paypal;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

import java.util.Date;
import java.util.List;

@Getter
@Setter
public class OrderResponse {

    private String id;

    private String intent;

    private String status;

    @JsonProperty("purchase_units")
    private List<PurchaseUnit> purchaseUnits;

    @JsonProperty("create_time")
    private Date createTime;

    private List<LinkDescription> links;
}
