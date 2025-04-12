package com.thanhtd.aerona.base.core;

import jodd.util.StringPool;
import org.apache.commons.lang3.math.NumberUtils;
import org.springframework.util.ObjectUtils;

import java.util.Arrays;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public final class AppUtils {
    private static final String EMAIL_REGEXP = "^[_A-Za-z0-9-]+(\\.[_A-Za-z0-9-]+)*@[A-Za-z0-9]+(\\.[A-Za-z0-9]+)*(\\.[A-Za-z]{2,})$";

    private AppUtils() {
        throw new IllegalStateException("Utility class");
    }

    public static boolean validateEmail(String email) {
        if (ObjectUtils.isEmpty(email)) return false;
        email = email.trim();
        if (email.startsWith("postmaster@") || email.startsWith("root@")) return false;

        Pattern pattern = Pattern.compile(EMAIL_REGEXP);
        Matcher matcher = pattern.matcher(email);
        if (matcher.matches()) {
            String localPath = email.substring(0, email.indexOf("@"));
            return localPath.length() >= 5 && localPath.length() <= 32;
        }
        return false;
    }

    public static boolean validatePhoneNumber(String phoneNumber) {
        String r = "^0\\d{9}$";
        return phoneNumber.matches(r);
    }

    public static Long[] parseIdsFromStr(String str) {
        return (Long[]) Arrays.stream(str.split(StringPool.DASH))
                .filter(NumberUtils::isCreatable)
                .map(Long::parseLong).toArray();
    }

    public static String generateUniqueId() {
        return UUID.randomUUID().toString();
    }
}
