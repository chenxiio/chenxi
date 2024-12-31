// #include <log4cxx/logger.h>
// #include <log4cxx/basicconfigurator.h>
// #include <log4cxx/helpers/exception.h>
// #include <log4cxx/patternlayout.h>
// #include <log4cxx/consoleappender.h>
// #include <log4cxx/fileappender.h>
//  using namespace log4cxx;
// using namespace log4cxx::helpers;
//  int main(int argc, char** argv)
// {
//     // 初始化日志配置
//     BasicConfigurator::configure();
//      // 创建日志记录器
//     LoggerPtr logger(Logger::getLogger("mylogger"));
//      // 创建输出格式
//     PatternLayoutPtr layout(new PatternLayout("%d{yyyy-MM-dd HH:mm:ss} [%t] %-5p %c %x - %m%n"));
//      // 创建输出目的地
//     FileAppenderPtr fileAppender(new FileAppender(layout, "mylog.log"));
//     ConsoleAppenderPtr consoleAppender(new ConsoleAppender(layout));
//      // 设置输出目的地的阈值
//     fileAppender->setThreshold(Level::getInfo());
//     consoleAppender->setThreshold(Level::getDebug());
//      // 添加输出目的地到日志记录器
//     logger->addAppender(fileAppender);
//     logger->addAppender(consoleAppender);
//      // 输出日志信息
//     LOG4CXX_INFO(logger, "Hello, log4cxx!");
//      return 0;
// }