#!/usr/bin/python3
import os
import logging
import time
import utilities
import requests
import json

filename = '/home/{}/sync/NetworkDeviceFinder.log'

try:
    name = utilities.get_user()
    filename = filename.format(name)
    os.remove(filename)
except OSError as error:
    pass


# Add the log message handler to the logger
logging.basicConfig(filename=filename,
                    format='%(asctime)s - %(levelname)s - %(message)s',
                    level=logging.INFO)
logging.info("Starting program")

class FileNotFound(Exception):
    '''Exception class for file checking'''

class Ping:

    def __init__(self):
        self.admin = ["192.168.1.101", "192.168.1.229"]
        self.send_data  = False
        self.server_address = ''

    def get_settings(self):
        '''Get config env var'''
        logging.info('get_settings()')
        name = utilities.get_user()
        config_name = '/home/{}/sync/config.json'
        config_name = config_name.format(name)
        try:
            if not os.path.isfile(config_name):
                raise FileNotFound('File is missing')
            with open(config_name) as file:
                data = json.load(file)
            self.server_address = '{}/alarm'.format(data["manager_address"])
            self.send_data = True
        except KeyError:
            logging.error("Variables not set")
        except FileNotFound:
            logging.error("File is missing")

    def ping_check(self, host):
        success = False
        try:
            response = os.system("ping -qc 3 {}".format(host))
            # and then check the response...
            if response == 0:
                logging.info('Success on checks')
                success = True
        except OSError as error:
            logging.error('Check failed on {}'.format(error))
        return success

    def publish_data(self, status):
        '''Send data to server if asked'''
        if self.send_data:
            state = True
            if status == 1:
                state = False
            data = {
                'status': state
            }
            try:
                response = requests.post(self.server_address, json=data, timeout=5)
                if response.status_code == 200:
                    logging.info("Requests successful")
                else:
                    logging.error('Response: {}'.format(response))
            except requests.ConnectionError as error:
                logging.error("Connection error: {}".format(error))
            except requests.Timeout as error:
                logging.error("Timeout on server: {}".format(error))

    def loop(self):
        '''Loop through sensor and publish'''
        while True:
            for admin in self.admin:
                self.publish_data(self.ping_check(admin))
                logging.info('Address: {}, Attempted ping: {}'.format(admin, True))
            time.sleep(60 * 3)


if __name__ == "__main__":
    ping = Ping()
    ping.get_settings()
    ping.loop()
