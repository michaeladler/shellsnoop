#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

static void print_usage(const char* prog)
{
    fprintf(stderr, "Usage: %s [-s SOCKET_FILE] PID\n", prog);
}

int main(int argc, char* argv[])
{
    char* socket_file = "/run/shellsnoop/shellsnoop.sock";
    int opt;
    while ((opt = getopt(argc, argv, "hs:")) != -1) {
        switch (opt) {
        case 's':
            socket_file = optarg;
            break;
        case 'h':
        default:
            print_usage(argv[0]);
            return EXIT_FAILURE;
        }
    }

    char* pid;
    if (optind < argc) {
        pid = argv[optind];
    } else {
        print_usage(argv[0]);
        return EXIT_FAILURE;
    }

    int fd;
    struct sockaddr_un address;
    char buf[4096];

    if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
        perror("Error creating socket");
        return EXIT_FAILURE;
    }

    memset(&address, 0, sizeof(address));
    address.sun_family = AF_UNIX;
    strncpy(address.sun_path, socket_file, sizeof(address.sun_path) - 1);

    if (connect(fd, (struct sockaddr*)&address, sizeof(address)) == -1) {
        perror("Error connecting to server");
        return EXIT_FAILURE;
    }

    if (send(fd, pid, strlen(pid), 0) == -1) {
        perror("Error sending message to server");
        return EXIT_FAILURE;
    }

    if (recv(fd, buf, sizeof(buf), 0) == -1) {
        perror("Error receiving message from server");
        return EXIT_FAILURE;
    }

    if (buf[0] != '\0') {
        puts(buf);
    }

    close(fd);

    return 0;
}
