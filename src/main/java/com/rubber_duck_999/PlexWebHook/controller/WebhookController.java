package com.rubber_duck_999.PlexWebHook.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.rubber_duck_999.PlexWebHook.model.Metadata;
import com.rubber_duck_999.PlexWebHook.service.NtfyService;

import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.http.ResponseEntity;
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
                if (event.equals("library.new")) {
                    Map<String, Object> metadataMap = (Map<String, Object>) payloadJSON.get("Metadata");
                    if (metadataMap != null) {
                        // Convert the metadataMap into a Metadata object using ObjectMapper
                        Metadata metadata = objectMapper.convertValue(metadataMap, Metadata.class);
                        // Use your notification logic here
                        StringBuilder sb = new StringBuilder();
                        sb.append("Type: ").append(metadata.getType()).append("\n");
                        sb.append("GrandparentTitle: ").append(metadata.getGrandparentTitle()).append("\n");
                        sb.append("ParentTitle: ").append(metadata.getParentTitle()).append("\n");
                        sb.append("Title: ").append(metadata.getTitle()).append("\n");
                        sb.append("Library Section: ").append(metadata.getLibrarySectionTitle()).append(" (")
                          .append(metadata.getLibrarySectionType()).append(")\n");
                        sb.append("Summary: ").append(metadata.getSummary()).append("\n");
                        System.out.println("Message: \n" + sb.toString());
                        ntfyService.sendNotification("webhook-events-7777", sb.toString());
                    }
                }
            }
        } catch (Exception e) {
            System.err.println("Failed to parse JSON payload: " + e.toString());
            System.out.println("Parsed Payload: " + payload.length());
        }
        return new ResponseEntity<>("Webhook received successfully", HttpStatus.OK);
    }
}
