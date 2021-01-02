# https://github.com/aws/aws-lambda-base-images/blob/python3.7/Dockerfile.python3.7
FROM scratch

ADD 0e313dfd7879f0a09a9f20da3168a11b265cab2030602938e5a21c8a89490dc4.tar.xz /
ADD 4ae7fb2df55c4579f90b43d4dcb7d2e6d2dac86e069469557230b1a8846e09fc.tar.xz /
ADD 59414281fed13668a189660fceefbec94e0fac6f6d2a6d3c3cfefd15da5a6438.tar.xz /
ADD 5fc256be16f8a1a54216b5d1102c3490d0505b6e0cb2a46c92bb3f602521881d.tar.xz /
ADD 91be3549ca8216eba4840177bbec5a1965439489276e0e8f8057ee0e1d953eaa.tar.xz /
ADD ab6a76ed98fa107b837951342d8f72e958259df5e80623c0a4780cca539bb3e9.tar.xz /

ENV LANG=en_US.UTF-8
ENV TZ=:/etc/localtime
ENV PATH=/var/lang/bin:/usr/local/bin:/usr/bin/:/bin:/opt/bin
ENV LD_LIBRARY_PATH=/var/lang/lib:/lib64:/usr/lib64:/var/runtime:/var/runtime/lib:/var/task:/var/task/lib:/opt/lib
ENV LAMBDA_TASK_ROOT=/var/task
ENV LAMBDA_RUNTIME_DIR=/var/runtime

WORKDIR /var/task

ENTRYPOINT ["/lambda-entrypoint.sh"]
