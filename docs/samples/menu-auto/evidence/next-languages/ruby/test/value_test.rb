require_relative '../lib/value'
raise 'unexpected value' unless Sample.value == 42
puts 'RUBY_SAMPLE_OK'
