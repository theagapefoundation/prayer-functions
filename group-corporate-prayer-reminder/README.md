# Group Corporate Prayer Reminder Function

The Prayer Reminder function is a dedicated component of the Prayer application, designed to send timely notifications to group members as reminders for their prayer schedules. It integrates seamlessly with the Prayer app, providing a reliable and efficient way to foster spiritual discipline and community engagement.

## Overview

This function is part of the broader suite of Google Cloud Functions supporting the [Prayer](https://github.com/theagapefoundation/prayer) application.

## Features

- **Automated Notifications:** Sends out reminder notifications to all members of a group who have opted in for this feature.
- **Customizable Schedule:** Allows users to set and manage their prayer reminder times according to their personal preferences and time zones.
- **Group Integration:** Works with the group feature of the Prayer app to ensure that all members receive reminders according to their group settings.

## Setup and Deployment

1. **Pre-requisites:**
   - Ensure you have access to the Google Cloud Platform.
   - Set up the Google Cloud SDK on your local machine.

2. **Configuration:**
   - Navigate to the Google Cloud Functions console.
   - Create a new function and specify the trigger and runtime settings.

3. **Deployment:**
   - Clone the repository containing the Prayer Reminder function.
   - Deploy the function to Google Cloud Functions using the Google Cloud SDK or through the Google Cloud Console.

## Usage

Once deployed, the Prayer Reminder function requires regular invocation to send notifications accurately. The function checks for scheduled notifications at the time of the call. Users can manage and view their reminder settings through the Prayer app interface.
