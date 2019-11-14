import socket
import sys, select

PORT = 1234
HEADER_SIZE = 10

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect((socket.gethostname(), PORT))

while True:
    read_stream, _, _ = select.select([sys.stdin, client], [], [])
    for stream in read_stream:
        if stream == client:
            full_msg = ''
            init_msg_received = False

            while True:
                msg = client.recv(16).decode('utf-8')

                if not init_msg_received:
                    msg_len = int(msg[:HEADER_SIZE])
                    init_msg_received = True
                
                full_msg += msg

                if len(full_msg) - HEADER_SIZE == msg_len:
                    print(f'Received: {full_msg[HEADER_SIZE:]}')
                    break
        else:
            msg = input()
            msg = f'{len(msg):<{HEADER_SIZE}}{msg}'
            print(f'Sending: {msg}')
            client.send(msg.encode('utf-8'))