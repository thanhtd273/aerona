package com.thanhtd.aerona.payment.dto.paypal;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class LinkDescription {

    private String href;

    private String rel;

    private String method;
}
