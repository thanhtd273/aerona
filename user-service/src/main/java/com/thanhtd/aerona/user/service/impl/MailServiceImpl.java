package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.service.MailService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.AppUtils;
import com.thanhtd.aerona.base.exception.LogicException;
import lombok.RequiredArgsConstructor;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

import java.util.Arrays;

@Service
@RequiredArgsConstructor
public class MailServiceImpl implements MailService {
    private final JavaMailSender mailSender;

    @Override
    public SimpleMailMessage sendSimpleMail(String to, String subject, String text) throws LogicException {
        // Validate destination email
        if (Boolean.FALSE.equals(AppUtils.validateEmail(to)))
            throw new LogicException(ErrorCode.INVALID_EMAIL);

        SimpleMailMessage message = new SimpleMailMessage();
        message.setText(text);
        message.setTo(to);
        message.setSubject(subject);
        mailSender.send(message);
        return message;
    }

    @Override
    public SimpleMailMessage sendSimpleMail(String[] to, String subject, String text) throws LogicException{
        // Validate destination emails
        if (Arrays.stream(to).anyMatch(email -> Boolean.FALSE.equals(AppUtils.validateEmail(email)))) {
            throw new LogicException(ErrorCode.INVALID_EMAIL);
        }
        SimpleMailMessage message = new SimpleMailMessage();
        message.setText(text);
        message.setTo(to);
        message.setSubject(subject);
        mailSender.send(message);
        return message;
    }
}
