package com.thanhtd.aerona.base.core;

import com.thanhtd.aerona.base.constant.ErrorCode;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class APIResponse {
    private Integer statusCode;

    private String statusMessage;

    private String description;

    private Object data;

    private Long took;

    public APIResponse(ErrorCode errorCode, String description, Long took, Object data) {
        this.statusCode = errorCode.getValue();
        this.statusMessage = errorCode.getMessage();
        this.description = description;
        this.data = data;
        this.took = took;
    }
}
