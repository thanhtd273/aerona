package com.thanhtd.aerona.user.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class Credential {
    private String email;

    private String password;
}
