package com.rubber_duck_999.PlexWebHook.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.rubber_duck_999.PlexWebHook.service.NtfyService;

import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.http.ResponseEntity;
import org.springframework.boot.autoconfigure.data.redis.RedisProperties.ClientType;
import org.springframework.http.HttpStatus;

import java.util.Map;

@RestController
@RequestMapping("/webhook")
public class WebhookController {

    private final NtfyService ntfyService;
    private final ObjectMapper objectMapper;

    public WebhookController(NtfyService ntfyService, ObjectMapper objectMapper) {
        this.ntfyService = ntfyService;
        this.objectMapper = objectMapper;
    }

    @PostMapping(consumes = "multipart/form-data")
    public ResponseEntity<String> handleMultipartWebhook(
            @RequestParam Map<String, String> formData,
            @RequestPart(value = "file", required = false ) MultipartFile file) {
        // Extract specific form fields from the payload
        String payload = formData.get("payload");
        try {
            @SuppressWarnings("unchecked")
            Map<String, Object> payloadJSON = objectMapper.readValue(payload, Map.class);
            // Get a specific key's value (e.g., "user")
            String event = (String) payloadJSON.get("event");
            if (event != null) {
                System.out.println("Event: " + event);
                String notificationMessage = "Event: " + event;
                ntfyService.sendNotification("webhook-events-6677", notificationMessage);
            }
        } catch (Exception e) {
            System.err.println("Failed to parse JSON payload");
            System.out.println("Parsed Payload: " + payload);
        }
        // Handle file upload if present
        if (file != null && !file.isEmpty()) {
            try {
                String fileContent = new String(file.getBytes());
                System.out.println("File content: " + fileContent);
            } catch (Exception e) {
                return new ResponseEntity<>("Failed to process file", HttpStatus.INTERNAL_SERVER_ERROR);
            }
        }

        return new ResponseEntity<>("Webhook received successfully", HttpStatus.OK);
    }
}
