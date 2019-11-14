import socket

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

try:
    sent = sock.sendto(
        'Hello from client'.encode('utf-8'),
        (socket.gethostname(), 1234)
    )
    data, server = sock.recvfrom(4096)
except Exception as e:
    print(f'Exception occured: {e}')
finally:
    sock.close()