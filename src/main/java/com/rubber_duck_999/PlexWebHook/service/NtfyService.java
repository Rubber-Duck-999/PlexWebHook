package com.rubber_duck_999.PlexWebHook.service;

import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;

@Service
public class NtfyService {

    private final RestTemplate restTemplate;

    public NtfyService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    public void sendNotification(String topic, String message) {
        String ntfyUrl = "https://ntfy.sh/" + topic;  // The ntfy server URL

        // Create the request headers and body
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.TEXT_PLAIN);

        HttpEntity<String> request = new HttpEntity<>(message, headers);

        // Send the notification
        restTemplate.postForEntity(ntfyUrl, request, String.class);
    }
}
