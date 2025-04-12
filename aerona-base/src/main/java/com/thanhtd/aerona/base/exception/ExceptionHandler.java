package com.thanhtd.aerona.base.exception;

import com.thanhtd.aerona.base.core.APIResponse;
import jakarta.servlet.http.HttpServletResponse;

public class ExceptionHandler {
    private ExceptionHandler() {
        throw new IllegalStateException("Static class");
    }

    public static APIResponse handleException(HttpServletResponse response, Exception e, long start) {
        if (e.getClass() == LogicException.class) {
            LogicException specEx = (LogicException) e;
            response.setStatus(specEx.getErrorCode().getValue());
            return new APIResponse(specEx.getErrorCode(), specEx.getMessage(), System.currentTimeMillis() - start, null);
        }
//        if (e.getClass() == BadCredentialsException.class) {
//            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
//            return new APIResponse(HttpServletResponse.SC_UNAUTHORIZED, "error", "Phone number or password is incorrect", System.currentTimeMillis() - start, null);
//        }
        response.setStatus(HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
        return new APIResponse(HttpServletResponse.SC_INTERNAL_SERVER_ERROR, "error", e.getMessage(), System.currentTimeMillis() - start, null );
    }
}
