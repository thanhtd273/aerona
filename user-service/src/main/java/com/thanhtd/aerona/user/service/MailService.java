package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.base.exception.LogicException;
import org.springframework.mail.SimpleMailMessage;

public interface MailService {

    SimpleMailMessage sendSimpleMail(String to, String subject, String text) throws LogicException;

    SimpleMailMessage sendSimpleMail(String[] to, String subject, String text) throws LogicException;
}
